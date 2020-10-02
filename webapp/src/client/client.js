// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';

import {id} from '../manifest';

export default class Client {
    setServerRoute(url) {
        this.url = url + '/plugins/' + id;
    }

    startWhiteboard = async (channelId, whiteboardId = 0) => {
        return doPost(`${this.url}/api/v1/whiteboards`, {channel_id: channelId, whiteboard_id: whiteboardId});
    }
}

export const doPost = async (url, body, headers = {}) => {
    const options = {
        method: 'post',
        body: JSON.stringify(body),
        headers,
    };

    const response = await fetch(url, Client4.getOptions(options));

    if (response.ok) {
        return response.json();
    }

    const text = await response.text();

    throw new ClientError(Client4.url, {
        message: text || '',
        status_code: response.status,
        url,
    });
};
