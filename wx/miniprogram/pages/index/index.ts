import { routing } from "../../utils/routing";

Page({
  isPageShowing: false,
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
      latitude: 23.099994,
      longitude: 113.324520,
    },
    scale: 10,
    markers: [
      {
        iconPath: "/resources/car.png",
        id: 0,
        latitude: 23.099994,
        longitude: 113.324520,
        width: 50,
        height: 50
      }, {
        iconPath: "/resources/car.png",
        id: 1,
        latitude: 23.099994,
        longitude: 114.324520,
        width: 50,
        height: 50
      }
    ]

  },
  onLoad(opt: Record<'car_id', string>) {
    const o: routing.LockOpts = opt
    console.log('unlocking car ', o.car_id);

  },
  onShow() {
    this.isPageShowing = true;
    const userInfo = getApp<IAppOption>().globalData.userInfo
    this.setData({
      avatarURL: userInfo?.avatarUrl
    })
  },
  onHide() {
    this.isPageShowing = false;
  },
  onScanClikced() {
    wx.scanCode({
      success: () => {
        //TODO: get car id from scan result
        const carID = 'car123'
        //作为参数值，需要转义
        const rediretcURL = routing.lock({ car_id: carID })
        wx.navigateTo({
          url: routing.register({
            redirectURL: rediretcURL
          })
        })
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
  moveCar() {
    const map = wx.createMapContext("map")
    const dest = {
      latitude: 23.099994,
      longitude: 113.324520,
    }

    const moveCar = () => {
      dest.latitude += 0.1
      dest.longitude += 0.1

      map.translateMarker({
        destination: {
          latitude: dest.latitude,
          longitude: dest.longitude,
        },
        markerId: 0,
        rotate: 0,
        duration: 5000,
        autoRotate: true,
        animationEnd: () => {
          if (this.isPageShowing) {
            moveCar()
          }
        }
      })
    }

    moveCar()

  }

})
