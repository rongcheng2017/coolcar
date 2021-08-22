import { routing } from "../../utils/routing"

// {{page}}.ts
Page({


  onRegisterTap(){
    wx.navigateTo({url:routing.register()})
  }
})