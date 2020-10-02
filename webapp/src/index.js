// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {id as pluginId} from './manifest';

import Icon from './components/icon';
import PostTypeWhiteboard from './components/post_type_whiteboard';
import {startWhiteboard} from './actions';
import Client from './client';

class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            (channel) => {
                startWhiteboard(channel.id)(store.dispatch, store.getState);
            },
            'Share Whiteboard',
        );
        registry.registerPostTypeComponent('custom_whiteboard', PostTypeWhiteboard);
        Client.setServerRoute(getServerRoute(store.getState()));
    }
}

window.registerPlugin(pluginId, new Plugin());

const getServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath;
};
