import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

// {{page}}.ts
Page({

  async onLoad(){
    const trips = await TripService.GetTrips(rental.v1.TripStatus.FINISHED)
  },
  onRegisterTap(){
    wx.navigateTo({url:routing.register()})
  }
})