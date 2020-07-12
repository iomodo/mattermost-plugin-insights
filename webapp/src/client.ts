import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import {id} from './manifest';

const apiUrl = `/plugins/${id}/api/v1`;

export const doGet = async (url: string) => {
    const {data} = await doFetchWithResponse(url, {method: 'get'});

    return data;
};


export const doFetchWithResponse = async (url: string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    let data;
    if (response.ok) {
        data = await response.json();

        return {
            response,
            data,
        };
    }

    data = await response.text();

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export function clientFetchTeams() {
    return doGet(`${apiUrl}/insights/teams`);
}

export function clientFetchChannels(teamID: string) {
    return doGet(`${apiUrl}/insights/channels?team_id=${teamID}&page=0&per_page=10`);
}

export function clientFetchPostData(teamID: string, channelID: string) {
    return doGet(`${apiUrl}/insights/post_data?team_id=${teamID}&channel_id=${channelID}`);
}

