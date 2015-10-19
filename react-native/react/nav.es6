'use strict'
/* @flow */

import Base from "./base"
import { connect } from 'react-redux'
import MetaNavigator from './router/meta-navigator'
import React from 'react'
import { StyleSheet } from 'react'
import Folders from './tabs/folders'
import Chat from './tabs/chat'
import People from './tabs/people'
//import Devices from './tabs/devices'
import NoTab from './tabs/no-tab'
import More from './tabs/more/index-desktop.es6'

import {FOLDER_TAB, CHAT_TAB, PEOPLE_TAB, DEVICES_TAB, MORE_TAB} from './constants/tabs'
import { switchTab } from './actions/tabbed-router'
import { Tab, Tabs, IconButton, Styles } from 'material-ui'
let { Colors, Typography } = Styles;

const tabToRootRouteParse = {
  [FOLDER_TAB]: Folders.parseRoute,
  [CHAT_TAB]: Chat.parseRoute,
  [PEOPLE_TAB]: People.parseRoute,
//  [DEVICES_TAB]: Devices.parseRoute,
  [MORE_TAB]: More.parseRoute
}

const menuItems = [
  { route: [FOLDER_TAB], text: 'Folders' },
  { route: [CHAT_TAB], text: 'Chat' },
  { route: [PEOPLE_TAB], text: 'People' }
]

export default class Nav extends Base {
  constructor(props) {
    super(props)
  }

  _renderContent (color, activeTab) {
    return (
      <div>
        {React.createElement(
          connect(state => {
            let elem = state.tabbedRouter.get('activeTab')
            // FIXME: After initial load, why do we end up with an array here?
            if (Array.isArray(elem)) {
              elem = elem[0]
            }
            return state.tabbedRouter.getIn(['tabs', elem]).toObject()
          })(MetaNavigator), {
            store: this.props.store,
            rootRouteParser: tabToRootRouteParse[activeTab] || NoTab.parseRoute
          }
        )}
      </div>
    )
  }

  _handleTabsChange (e, key, payload) {
    console.log('should switch to ' + e)
    this.props.dispatch(switchTab(e))
  }

  render () {
    const {dispatch} = this.props
    const activeTab = this.props.tabbedRouter.get('activeTab')

    let styles = {
      div: {
        position: 'absolute',
        left: 48,
        backgroundColor: Colors.cyan500,
        width: 0,
        height: 48,
      },
      headline: {
        fontSize: 24,
        lineHeight: '32px',
        paddingTop: 16,
        marginBottom: 12,
        letterSpacing: 0,
        fontWeight: Typography.fontWeightNormal,
        color: Typography.textDarkBlack,
      },
      iconButton: {
        position: 'absolute',
        left: 0,
        backgroundColor: Colors.cyan500,
        color: 'white',
        marginRight: 0,
      },
      iconStyle: {
        color: Colors.white,
      },
      tabs: {
        position: 'relative',
      },
      tabsContainer: {
        position: 'relative',
        paddingLeft: 0,
        width: '70%'
      },
    };

    return (
      <div style={styles.tabsContainer}>
        <Tabs valueLink={{value: activeTab, requestChange: this._handleTabsChange.bind(this)}}>
          <Tab label="More" value={MORE_TAB} >
            {this._renderContent('#aaaaaa', activeTab)}
          </Tab>
          <Tab label="Folders" value={FOLDER_TAB} >
            {this._renderContent('#aaaaaa', activeTab)}
          </Tab>
          <Tab label="Chat" value={CHAT_TAB}>
            {this._renderContent('#aaaaaa', activeTab)}
          </Tab>
          <Tab label="People" value={PEOPLE_TAB}>
            {this._renderContent('#aaaaaa', activeTab)}
          </Tab>
        </Tabs>
      </div>
    )
  }
}
