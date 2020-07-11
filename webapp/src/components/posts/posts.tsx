import React, {useEffect, useState} from 'react';
import ItemSelector from '../item_selector/item_selector';
import {Team, Channel} from '../../types/posts';
import {clientFetchTeams, clientFetchChannels} from '../../client';

const Posts = () => {
    const [teamId, setTeamId] = useState<string>();
    const [channelId, setChannelId] = useState<string>();

    useEffect(() => {
        console.log('teamId', teamId, 'channelId', channelId);
    }, [teamId, channelId]);

    return(
        <div>
            <h1>Post Statistics</h1>
            Teams
            <ItemSelector
                getItems={clientFetchTeams}
                argument={''}
                onSelectedChange={(teamId) => {setTeamId(teamId);}}
            />
            Channels
            <ItemSelector
                getItems={clientFetchChannels}
                argument={teamId?teamId:''}
                onSelectedChange={(channelId) => {setChannelId(channelId);}}
            />
        </div>
        
    );
}

async function fetchChannels() {
    const teams: Channel[] = [{
        id: '3',
        display_name: 'chan1',
        name: 'chan',
    }, {
        id: '24',
        display_name: 'chan2',
        name: 'chan_chan',
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
