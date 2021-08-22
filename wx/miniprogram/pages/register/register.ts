// {{page}}.ts
Page({
  rediretcURL: '',
  data: {
    licNo: '',
    name: '',
    genderIndex: 0,
    genders: ['未知', '男', '女', '其他'],
    birthDate: '1990-01-01',
    licImgURL: '',
    state: 'UNSUBMITTED' as 'UNSUBMITTED' | 'PENDING' | 'VERIFIED',
  },

  onLoad(opt) {
    if (opt.redirect) {
      this.rediretcURL = opt.redirect
    }
  },
  onUploadLic() {
    wx.chooseImage({
      count: 1,
      success: res => {
        this.setData({
          licImgURL: res.tempFilePaths[0]
        })
        //TODO: upload image
        setTimeout(() => {
          this.setData({
            licNo: '32423444',
            name: '张三',
            genderIndex: 1,
            birthDate: '1990-10-10'
          })
        }, 1000)
      }
    })
  }
  ,
  onGenderChange(e: any) {
    this.setData({
      genderIndex: e.detail.value
    })
  },
  onBrithDateChange(e: any) {
    this.setData({
      birthDate: e.detail.value
    })
  },
  onSubmit() {
    //TODO: submit the form to server
    this.setData({
      state: 'PENDING'
    })
    setTimeout(() => {
      this.onLicVerified()
    }, 3000);
  },
  onReSubmit() {
    this.setData({
      state: 'UNSUBMITTED',
      licImgURL: '',
    })
  },
  onLicVerified() {
    this.setData({
      state: 'VERIFIED'
    })
    if (this.rediretcURL) {
      wx.redirectTo({ url: decodeURIComponent(this.rediretcURL) })
    }

  }

})