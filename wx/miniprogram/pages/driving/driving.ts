import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

const upadteIntervalSec = 5

function formatDuration(sec: number) {
  const padString = (n: number) => n < 10 ? '0' + n.toFixed(0) : n.toFixed(0)

  const h = Math.floor(sec / 3600)
  sec -= 3600 * h
  const m = Math.floor(sec / 60)
  sec -= 60 * m
  const s = Math.floor(sec)

  return `${padString(h)}:${padString(m)}:${padString(s)}`

}

function formatFee(cents: number) {
  return (cents / 100).toFixed(2)
}

Page({

  timer: undefined as undefined | number,
  tripID: '',
  data: {
    location: {
      latitude: 32.92,
      longitude: 118.46,
    },
    scale: 14,
    elapsed: '00:00:00',
    fee: '0.00'
  },
  onLoad(opt: Record<'trip_id', string>) {
    const o: routing.DrivingOpts = opt
    console.log('current trip', o.trip_id);
    this.tripID = o.trip_id
    TripService.GetTrip(o.trip_id).then(console.log)
    this.setupLocationUpdator()
    this.setupTimer(this.tripID)
  },
  onUnload() {
    wx.stopLocationUpdate()
    if (this.timer) {
      clearInterval(this.timer)
    }
  },
  setupLocationUpdator() {
    wx.startLocationUpdate({
      fail: console.error
    })
    wx.onLocationChange(loc => {
      console.log(loc)
      this.setData({
        location: {
          latitude: loc.latitude,
          longitude: loc.longitude,
        }
      })
    })
  },
  async setupTimer(tripID: string) {
    const trip = await TripService.GetTrip(tripID)
    if (trip.status !== rental.v1.TripStatus.IN_PROGRESS) {
      console.log('trip not in progress');
      return
    }
    let secSinceLastUpdate = 0
    let lastUpdateDurationSec = trip.current!.timestampSec as number
    this.setData({
      elapsed: formatDuration(lastUpdateDurationSec),
      fee: formatFee(trip.current!.feeCent!)
    })
    this.timer = setInterval(() => {
      secSinceLastUpdate++
      //5s refresh 
      if (secSinceLastUpdate % upadteIntervalSec == 0) {
        TripService.GetTrip(tripID).then(trip => {
          lastUpdateDurationSec = trip.current!.timestampSec!
          secSinceLastUpdate = 0
          this.setData({
            fee: formatDuration(trip.current!.feeCent!)
          })
        }).catch(console.error)
      }
      this.setData({
        elapsed: formatDuration(lastUpdateDurationSec + secSinceLastUpdate),
      })
    }, 1000)
  },
  onEndTripTap() {
    TripService.finishTrip(this.tripID).then(() => {
      wx.redirectTo({
        url: routing.mytrips(),
      })
    }).catch(err => {
      console.error(err)
      wx.showToast({
        title: '结束行程失败',
        icon: 'none',
      })
    })
  }
})