import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { padString } from "../../utils/format"
import { routing } from "../../utils/routing"

function formatDate(mills: number) {
  const dt = new Date(mills)
  const y = dt.getFullYear()
  const m = dt.getMonth() + 1
  const d = dt.getDate()
  return `${padString(y)}-${padString(m)}-${padString(d)}`

}
// {{page}}.ts
Page({
  rediretcURL: '',
  profileRefresher: 0,
  data: {
    licNo: '',
    name: '',
    genderIndex: 0,
    genders: ['未知', '男', '女'],
    birthDate: '1990-01-01',
    licImgURL: '',
    state: rental.v1.IdentityStatus[rental.v1.IdentityStatus.UNSUBMITTED],
  },
  renderProfile(p: rental.v1.IProfile) {
    this.setData(
      {
        licNo: p.identity?.licNumber || '',
        name: p.identity?.name || '',
        genderIndex: p.identity?.gender || 0,
        birthDate: formatDate(p.identity?.birthDateMillis || 0),
        state: rental.v1.IdentityStatus[p.identityStatus || 0],
      }
    )
  },
  onLoad(opt: Record<'redirect', string>) {
    const o: routing.RegisterOpts = opt
    if (o.redirect) {
      this.rediretcURL = decodeURIComponent(o.redirect)
    }
    ProfileService.getProfile().then(p => {
      this.renderProfile(p)
    })
  },

  onUnload() {
    this.clearProfileRefresher()
  },
  onUploadLic() {
    wx.chooseImage({
      count: 1,
      success: res => {
        this.setData({
          licImgURL: res.tempFilePaths[0]
        })
        const data=wx.getFileSystemManager().readFileSync(res.tempFilePaths[0])

        wx.request({
          method:'PUT',
          url:'https://coolcar-1307431695.cos.ap-beijing.myqcloud.com/account_1%2F613f2ed1642df9678f45c188?q-sign-algorithm=sha1&q-ak=AKIDVpseIflbCYeT2KL8gxxg8KFtHPGq9CyB&q-sign-time=1631530705%3B1631531705&q-key-time=1631530705%3B1631531705&q-header-list=host&q-url-param-list=&q-signature=dca8075021d818168e7e7ce9494d0839c950645a',
          data,
          success:console.log,
          fail:console.error
        })
      }
    })
  }
  ,
  onGenderChange(e: any) {
    this.setData({
      genderIndex: parseInt(e.detail.value)
    })
  },
  onBrithDateChange(e: any) {
    this.setData({
      birthDate: e.detail.value
    })
  },
  onSubmit() {
    ProfileService.submitProfile({
      licNumber: this.data.licNo,
      name: this.data.name,
      gender: this.data.genderIndex,
      birthDateMillis: Date.parse(this.data.birthDate)
    }).then(p => {
      this.renderProfile(p)
      this.scheduleProfileRefresher()
    })
  },
  scheduleProfileRefresher() {
    this.profileRefresher = setInterval(() => {
      ProfileService.getProfile().then(p => {
        this.renderProfile(p)
        if (p.identityStatus !== rental.v1.IdentityStatus.PENDING) {
          this.clearProfileRefresher()
        }
        if (p.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
          this.onLicVerified()
        }
      })
    }, 1000)
  },
  clearProfileRefresher() {
    if (this.profileRefresher) {
      clearInterval(this.profileRefresher)
      this.profileRefresher = 0
    }
  },
  onReSubmit() {
    ProfileService.clearProflie().then(p => this.renderProfile(p))
  },
  onLicVerified() {
    if (this.rediretcURL) {
      wx.redirectTo({ url: this.rediretcURL })
    }
  }

})