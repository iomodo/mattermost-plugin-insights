import React, {useEffect, useState} from 'react';
import ItemSelector from '../item_selector/item_selector';
import {Team, Channel, Item} from '../../types/posts';
import {clientFetchTeams, clientFetchChannels, clientFetchPostData} from '../../client';
import {
    LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend,
  } from 'recharts';

const Posts = () => {
    const [teamId, setTeamId] = useState<string>();
    const [channelId, setChannelId] = useState<string>();
    const [fetchChannels, setFetchChannels] = useState<() => Promise<Item[]>>();
    const [postData, setPostData] = useState<[]>();

    useEffect(() => {
        console.log('teamId', teamId, 'channelId', channelId);
        setFetchChannels(() => () => {
            return clientFetchChannels(teamId?teamId:'');
        })
    }, [teamId]);

    useEffect(() => {
        console.log('trying to fetch data');
        clientFetchPostData(teamId?teamId:'', channelId?channelId:'').then(data => {
            console.log(data)
            setPostData(data);
        })
    }, [teamId, channelId]);

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
            <LineChart
                width={500}
                height={200}
                data={postData}
                margin={{
                    top: 10, right: 30, left: 0, bottom: 0,
                }}
            >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Line connectNulls type="monotone" dataKey="value" stroke="#8884d8" fill="#8884d8" />
            </LineChart>
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
