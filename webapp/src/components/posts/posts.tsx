import React from 'react';
import TeamSelector from '../team_selector/team_selector'
import ItemSelector from '../item_selector/item_selector';
import {Team} from '../../types/posts';

const Posts = () => {
    return(
        <div>
            <h1>Post Statistics</h1>
            Teams
            <ItemSelector
                getItems={fetchTeams}
                onSelectedChange={(teamId) => {console.log(teamId);}}
            />
        </div>
        
    );
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
