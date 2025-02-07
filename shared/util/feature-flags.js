/* @flow */

type FeatureFlag = {
  tracker: 'v1' | 'v2'
}
const ff: FeatureFlag = {
  tracker: process.env.KEYBASE_TRACKER_V2 ? 'v2' : 'v1'
}
export default ff
