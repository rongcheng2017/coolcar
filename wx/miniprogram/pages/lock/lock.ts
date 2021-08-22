const shareLocationKey = 'share_location'

Page({
    data: {
        avatarURL: '',
        shareLocation: false,
    },
    async onLoad() {
        const userInfo =getApp<IAppOption>().globalData.userInfo
        this.setData({
            shareLocation: wx.getStorageSync(shareLocationKey) || false,
            avatarURL:userInfo?.avatarUrl
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


    }
})