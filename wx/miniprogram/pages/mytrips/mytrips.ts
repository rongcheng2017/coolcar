import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"
import { getUserInfo } from "../../utils/wxapi"

// {{page}}.ts
const licStatusMap = new Map([
  [rental.v1.IdentityStatus.UNSUBMITTED, '未认证'],
  [rental.v1.IdentityStatus.PENDING, '未认证'],
  [rental.v1.IdentityStatus.VERIFIED, '已认证'],
])
Page({
  data: {
    licStatus: licStatusMap.get(rental.v1.IdentityStatus.UNSUBMITTED),
  },
  async onLoad() {
    const trips = await TripService.GetTrips(rental.v1.TripStatus.FINISHED)

  },
  onShow() {
    ProfileService.getProfile().then(p => {
      this.setData({
        licStatus: licStatusMap.get(p.identityStatus || 0)
      })
    })
  },
  onRegisterTap() {
    wx.navigateTo({ url: routing.register() })
  },
  onGetUserInfo() {
    getUserInfo().then(res => {
      console.log(res)
    })
  }
})