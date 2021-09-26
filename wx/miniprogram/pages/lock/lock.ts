import { IAppOption } from "../../appoption";
import { CarService } from "../../service/car";
import { car } from "../../service/proto_gen/car/car_pb";
import { TripService } from "../../service/trip";
import { routing } from "../../utils/routing";

const shareLocationKey = 'share_location'

Page({
    carID: '',
    carRefresher: 0,
    data: {
        avatarURL: '',
        shareLocation: true,
    },
    async onLoad(opt: Record<'car_id', string>) {
        const o: routing.LockOpts = opt
        this.carID = o.car_id
        const userInfo = getApp<IAppOption>().globalData.userInfo
        this.setData({
            shareLocation: wx.getStorageSync(shareLocationKey) || true,
            avatarURL: userInfo?.avatarUrl
        })
    },
    onGetUserInfo(e: any) {
        const userInfo: WechatMiniprogram.UserInfo = e.detail.userInfo
        //处理拒接的情况
        if (userInfo) {
            this.setData({
                avatarURL: userInfo.avatarUrl
            })
            getApp<IAppOption>().globalData.userInfo = userInfo
        }

    },
    onShareLocation(e: any) {
        this.data.shareLocation = e.detail.value
        wx.setStorageSync(shareLocationKey, this.data.shareLocation)
    },
    onUnlockTap() {
        wx.getLocation({
            type: 'gcj02',
            success: async loc => {

                const location = {
                    latitude: loc.latitude,
                    longitude: loc.longitude,
                }

                if (!this.carID) {
                    console.error('no carID specified');
                    return
                }
                const trip = await TripService.CreateTrip({
                    start: location,
                    carId: this.carID,
                    avatarUrl: this.data.shareLocation ? this.data.avatarURL : '',
                })
                if (!trip.id) {
                    console.error('no tripID in response', trip);
                    return
                }
                wx.showLoading({ title: '开锁中', mask: true, })

                this.carRefresher = setInterval(async () => {
                    const c = await CarService.getCar(this.carID)
                    if (c.status === car.v1.CarStatus.UNLOCKED) {
                        this.clearCarRefresher()
                        wx.redirectTo({
                            url: routing.driving({
                                trip_id: trip.id
                            }),
                            complete: () => {
                                wx.hideLoading()
                            }
                        })
                    }
                }, 2000)
            },
            fail: () => {
                wx.showToast({
                    icon: 'none',
                    title: '请前往设置页面授权位置信息'
                })
            }

        })

    },
    onUnload() {
        this.clearCarRefresher()
    },
    clearCarRefresher() {
        if (this.carRefresher) {
            clearInterval(this.carRefresher)
            this.carRefresher = 0
        }
    }
})