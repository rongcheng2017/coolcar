import { rental } from "./proto_gen/rental/rental_pb";
import { Coolcar } from "./request";

export namespace TripService {
    export function CreateTrip(req: rental.v1.ICreateTripRequest): Promise<rental.v1.TripEntity> {
        return Coolcar.sendRequesWithAutyRetry({
            method: 'POST',
            path: '/v1/trip',
            data: req,
            respMarshaller: rental.v1.TripEntity.fromObject
        })
    }
    export function GetTrip(id: string): Promise<rental.v1.ITrip> {
        return Coolcar.sendRequesWithAutyRetry(
            {
                method: 'GET',
                path: `/v1/trip/${encodeURIComponent(id)}`,
                respMarshaller: rental.v1.Trip.fromObject,
            }
        )
    }

    export function GetTrips(s?: rental.v1.TripStatus): Promise<rental.v1.IGetTripsResponse> {
        let path = '/v1/trips'
        if (s) {
            path += `?status=${s}`
        }
        return Coolcar.sendRequesWithAutyRetry({
            method: 'GET',
            path,
            respMarshaller: rental.v1.GetTripsResponse.fromObject
        })

    }

    export function finishTrip(id: string) {
        return updateTrip({ id, endTrip: true })

    }

    function updateTrip(r: rental.v1.IUpdateTripRequest): Promise<rental.v1.ITrip> {
        if (!r.id) {
            return Promise.reject("must specify id")
        }
        return Coolcar.sendRequesWithAutyRetry({
            method: 'PUT',
            path: `/v1/trip/${encodeURIComponent(r.id)}`,
            data: r,
            respMarshaller: rental.v1.Trip.fromObject,
        })

    }
}