const shareLocationKey = 'share_location'

Page({
    data: {
        avatarURL: '',
        shareLocation: true,
    },
    async onLoad(opt) {
        console.log('unlocking car ',opt.car_id);
        
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
            success: loc => {
                console.log('starting a trip', {
                    location: {
                        latitude: loc.latitude,
                        longitude: loc.longitude,
                    },
                    avatarURL: this.data.shareLocation ? this.data.avatarURL : '',
                    // carID:'33322'

                })
                const tripID = 'trip456'
                wx.showLoading({ title: '开锁中', mask: true })
                setTimeout(() => {
                    wx.redirectTo({
                        url: `/pages/driving/driving?trip_id=${tripID}`,
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