import { IAppOption } from "../../appoption";
import { CarService } from "../../service/car";
import { ProfileService } from "../../service/profile";
import { rental } from "../../service/proto_gen/rental/rental_pb";
import { TripService } from "../../service/trip";
import { routing } from "../../utils/routing";

interface Marker {
  iconPath: string
  id: number
  latitude: number
  longitude: number
  width: number
  height: number
}

const defaultAvatar = '/resources/car.png'
const initialLat = 29.76090514382441
const initialLng = 121.86952988416326

Page({
  isPageShowing: false,
  socket: undefined as WechatMiniprogram.SocketTask | undefined,
  data: {
    avatarURL: '',
    setting: {
      skew: 0,
      rotate: 0,
      showLocation: true,
      showScale: true,
      subKey: '',
      layerStyle: -1,
      enableZoom: true,
      enableScroll: true,
      enableRotate: false,
      showCompass: false,
      enable3D: false,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,
    },
    location: {
      latitude: initialLat,
      longitude: initialLng,
    },
    scale: 16,
    markers: [] as Marker[]

  },
  onShow() {
    this.isPageShowing = true;
    const userInfo = getApp<IAppOption>().globalData.userInfo
    this.setData({
      avatarURL: userInfo?.avatarUrl
    })
    if (!this.socket) {
      this.setData({
        markers: []
      },()=>this.setupCarPosUpdater())
    }
  },
  onHide() {
    this.isPageShowing = false;
    if (this.socket) {
      this.socket.close({
        success: () => {
          this.socket = undefined
        }
      })
    }
  },
  async onScanTap() {
    //扫描之前先确认有没有正在使用的行程
    const trips = await TripService.GetTrips(rental.v1.TripStatus.IN_PROGRESS)
    if ((trips.trips?.length || 0) > 0) {
      // await this.selectComponent('#tripModal').showModal()
      wx.navigateTo({ url: routing.driving({ trip_id: trips.trips![0].id! }) })
      return
    }
    wx.scanCode({
      success: async () => {
        const carID = '615034936235ac04a7aca3c9'
        //作为参数值，需要转义
        const lockURL = routing.lock({ car_id: carID })
        const prof = await ProfileService.getProfile()
        if (prof.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
          wx.navigateTo({
            url: lockURL
          })
        } else {
          // await this.selectComponent('#licModal').showModal()
          wx.navigateTo({
            url: routing.register({
              redirectURL: lockURL
            })
          })
        }
      },
      fail: console.error
    })
  },
  onMyLocationTap() {
    wx.getLocation(
      {
        type: 'gcj02',
        success: res => {
          this.setData({
            location: {
              latitude: res.latitude,
              longitude: res.longitude,
            },
          })
        },
        fail: () => {
          wx.showToast({
            title: '请前往设置页授权',
            icon: 'none'
          })
        }
      }
    )
  },
  onMyTripsTap() {
    wx.navigateTo({ url: routing.mytrips() })
  },
  setupCarPosUpdater() {
    const markersByCarID = new Map<string, Marker>()
    const map = wx.createMapContext("map")
    let translationInProgress = false
    const endTranslation = () => {
      translationInProgress = false

    }
    this.socket = CarService.subscribe(car => {
      if (!car.id || translationInProgress || !this.isPageShowing) {
        return
      }
      const newLat = car.car?.position?.latitude || initialLat
      const newLng = car.car?.position?.longitude || initialLng
      const marker = markersByCarID.get(car.id)

      if (!marker) {
        const newMarker: Marker = {
          id: this.data.markers.length,
          iconPath: car.car?.driver?.avatarUrl || defaultAvatar,
          latitude: newLat,
          longitude: newLng,
          height: 20,
          width: 20,
        }
        markersByCarID.set(car.id, newMarker)
        this.data.markers.push(newMarker)
        translationInProgress = true
        this.setData({
          markers: this.data.markers,
        }, endTranslation)
        return
      }

      const newAvatar = car.car?.driver?.avatarUrl || defaultAvatar
      if (marker.iconPath !== newAvatar) {
        marker.iconPath = newAvatar
        marker.latitude = newLat
        marker.longitude = newLng
        translationInProgress = true
        this.setData({
          markers: this.data.markers
        }, endTranslation)
        return
      }

      if (marker.latitude !== newLat || marker.longitude !== newLng) {
        translationInProgress = true
        // Move cars.
        map.translateMarker({
          markerId: marker.id,
          destination: {
            latitude: newLat,
            longitude: newLng,
          },
          autoRotate: false,
          rotate: 0,
          duration: 90,
          animationEnd: endTranslation,
        })
      }
    })
  }
  ,

})
