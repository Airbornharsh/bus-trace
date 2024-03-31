export interface Bus {
  ID: string
  name: string
  lat: number
  long: number
}

interface busClose {
  close: boolean
}

enum BusReqWhich {
  BusData = 'busData',
  BusClose = 'busClose'
}

export interface BusReq {
  busId: string
  busData: busData
  busClose: busClose
  which: BusReqWhich
}
interface busData {
  lat: number
  long: number
  index: number
}

interface busMessage {
  busId: string
  message: string
}

interface busUserLocation {
  userId: string
  lat: number
  long: number
}

enum BusResWhich {
  BusData = 'busData',
  BusMessage = 'busMessage',
  BusUserList = 'busUserList',
  BusUserLocation = 'busUserLocation'
}

export interface BusRes {
  busId: string
  busData: busData
  busMessage: busMessage
  busUserList: string[]
  busUserLocation: busUserLocation
  which: BusResWhich
}
