import camelcaseKeys = require("camelcase-keys");
import { IAppOption } from "./appoption"
import { coolcar } from "./service/proto_gen/trip_pb"

// app.ts
App<IAppOption>({
  globalData: {},
  onLaunch() {
    wx.request({
      url: "http://localhost:8080/trip/trip123",
      method: 'GET',
      success: res => {
        const getTripRes = coolcar.GetTripResponse.fromObject(camelcaseKeys(res.data as object, { deep: true }))
        console.log(getTripRes);

        //enum -> string
        console.log(coolcar.TripStatus[getTripRes.trip?.status!])

      },
      fail: console.error,
    })
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    // wx.setStorageSync('logs', logs)

    // 登录
    wx.login({
      success: res => {
        console.log(res.code)
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
    })
  },

})

function camelcaseKey(): { [k: string]: any; } {
  throw new Error("Function not implemented.");
}
