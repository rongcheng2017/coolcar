import { IAppOption } from "./appoption"
import { Coolcar } from "./service/request";

// app.ts
App<IAppOption>({
  globalData: {},
  onLaunch() {
    // 登录
    Coolcar.login()
  }
})
