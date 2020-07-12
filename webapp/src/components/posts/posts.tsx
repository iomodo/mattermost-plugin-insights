import React, {useEffect, useState} from 'react';
import ItemSelector from '../item_selector/item_selector';
import {Team, Channel, Item} from '../../types/posts';
import {clientFetchTeams, clientFetchChannels} from '../../client';

const Posts = () => {
    const [teamId, setTeamId] = useState<string>();
    const [channelId, setChannelId] = useState<string>();
    const [fetchChannels, setFetchChannels] = useState<() => Promise<Item[]>>();

    useEffect(() => {
        console.log('teamId', teamId, 'channelId', channelId);
        setFetchChannels(() => () => {
            return clientFetchChannels(teamId?teamId:'');
        })
    }, [teamId]);

    return(
        <div>
            <h1>Post Statistics</h1>
            Teams
            <ItemSelector
                getItems={clientFetchTeams}
                onSelectedChange={(teamId) => {setTeamId(teamId);}}
            />
            Channels
            <ItemSelector
                getItems={fetchChannels?fetchChannels:defaultFetchChannels}
                onSelectedChange={(channelId) => {setChannelId(channelId);}}
            />
        </div>
        
    );
}

async function defaultFetchChannels() {
    const teams: Item[] = [{
        id: '3',
        display_name: 'chan1',
    }, {
        id: '24',
        display_name: 'chan2',
    }]
    return teams;
}

async function fetchTeams() {
    const teams: Team[] = [{
        id: '1',
        display_name: 'bla1',
        name: 'bla_bla',
    }, {
        id: '2',
        display_name: 'bla2',
        name: 'bla_bla',
    }]
    return teams;
}

export default Posts;
