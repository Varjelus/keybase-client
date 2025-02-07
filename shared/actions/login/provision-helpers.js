import * as Constants from '../../constants/login'
import QRCodeGen from 'qrcode-generator'

export function defaultModeForDeviceRoles (myDeviceRole, otherDeviceRole, brokenMode) {
  switch (myDeviceRole + otherDeviceRole) {
    case Constants.codePageDeviceRoleExistingComputer + Constants.codePageDeviceRoleNewComputer:
      return Constants.codePageModeEnterText
    case Constants.codePageDeviceRoleNewComputer + Constants.codePageDeviceRoleExistingComputer:
      return Constants.codePageModeShowText

    case Constants.codePageDeviceRoleExistingComputer + Constants.codePageDeviceRoleNewPhone:
      return Constants.codePageModeShowCodeOrEnterText
    case Constants.codePageDeviceRoleNewPhone + Constants.codePageDeviceRoleExistingComputer:
      return Constants.codePageModeScanCode

    case Constants.codePageDeviceRoleExistingPhone + Constants.codePageDeviceRoleNewComputer:
      return Constants.codePageModeScanCode
    case Constants.codePageDeviceRoleNewComputer + Constants.codePageDeviceRoleExistingPhone:
      return Constants.codePageModeShowCodeOrEnterText

    case Constants.codePageDeviceRoleExistingPhone + Constants.codePageDeviceRoleNewPhone:
      return brokenMode ? Constants.codePageModeShowText : Constants.codePageModeShowCode
    case Constants.codePageDeviceRoleNewPhone + Constants.codePageDeviceRoleExistingPhone:
      return brokenMode ? Constants.codePageModeEnterText : Constants.codePageModeScanCode
  }
  return null
}

export function qrGenerate (code) {
  const qr = QRCodeGen(4, 'L')
  qr.addData(code)
  qr.make()
  let tag = qr.createImgTag(10)
  const src = tag.split(' ')[1]
  const qrCode = src.split('\"')[1]
  return qrCode
}
