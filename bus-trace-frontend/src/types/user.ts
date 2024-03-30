import { Location } from './other'

export interface User {
  ID: string
  name: string
  email: string
  phone: string
  busOwner: boolean
  busId: string | null
}

enum UserReqWhich {
  UserLocation = 'userLocation'
}

export interface UserReq {
  userLocation: Location
  which: UserReqWhich
}

enum UserResWhich {
  UserMessage = 'userMessage'
}

interface userBusData {
  busId: string
  lat: number
  long: number
}

interface userMessage {
  message: string
}

export interface UserRes {
  userId: string
  userBusData: userBusData
  userMessage: userMessage
  which: UserResWhich
}
