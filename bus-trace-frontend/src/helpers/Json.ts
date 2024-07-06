export const validate = (data: string) => {
  try {
    JSON.parse(data)
    return true
  } catch (err) {
    return false
  }
}
