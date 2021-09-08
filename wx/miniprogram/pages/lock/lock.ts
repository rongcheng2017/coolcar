import { IAppOption } from "../../appoption";
import { TripService } from "../../service/trip";
import { routing } from "../../utils/routing";

const shareLocationKey = 'share_location'

Page({
    carID: '',
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
        const shareLocation: boolean = e.detail.value
        wx.setStorageSync(shareLocationKey, shareLocation)


    },
    onUnlockTap() {
        wx.getLocation({
            type: 'gcj02',
            success: async loc => {

                const location = {
                    latitude: loc.latitude,
                    longitude: loc.longitude,
                }
                console.log('starting a trip', {
                    location: location,
                    avatarURL: this.data.shareLocation ? this.data.avatarURL : '',
                })
                if (!this.carID) {
                    console.error('no carID specified');
                    return
                }
                const trip = await TripService.CreateTrip({
                    start: location,
                    carId: this.carID,
                })
                if (!trip.id) {
                    console.error('no tripID in response', trip);
                    return
                }
                wx.showLoading({ title: '开锁中', mask: true, })
                setTimeout(() => {
                    wx.redirectTo({
                        url: routing.driving({
                            trip_id: trip.id
                        }),
                        complete: () => {
                            wx.hideLoading()
                        }
                    })
                }, 2000)

            },
            fail: () => {
                wx.showToast({
                    icon: 'none',
                    title: '请前往设置页面授权位置信息'
                })
            }

        })

    }
})