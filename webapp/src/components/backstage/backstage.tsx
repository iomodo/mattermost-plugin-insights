// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import classNames from 'classnames';

import './backstage.scss';
import {BackstageArea} from '../../types/backstage';
import Posts from '../posts/posts';
import Users from '../users/users';

interface Props {
    onBack: () => void;
    selectedArea: BackstageArea;
    setSelectedArea: (area: BackstageArea) => void;
}

const Backstage = ({onBack, selectedArea, setSelectedArea}: Props): React.ReactElement<Props> => {
    console.log("selectedArea", selectedArea)
    let activeArea = <Posts/>;
    if (selectedArea === BackstageArea.Users) {
        activeArea = <Users/>;
    }
    console.log("activeArea", activeArea)

    return (
        <div className='Backstage'>
            <div className='Backstage__sidebar'>
                <div className='Backstage__sidebar__header'>
                    <div
                        className='cursor--pointer'
                        onClick={onBack}
                    >
                        <i className='icon-arrow-left mr-2 back-icon'/>
                        {'Back to Mattermost'}
                    </div>
                </div>
                <div className='menu'>
                    <div 
                        className={classNames('menu-title', {active: selectedArea === BackstageArea.Dashboard})}
                        onClick={() => setSelectedArea(BackstageArea.Dashboard)}
                    >
                        {'Dashboard'}
                    </div>
                    <div
                        className={classNames('menu-title', {active: selectedArea === BackstageArea.Posts})}
                        onClick={() => setSelectedArea(BackstageArea.Posts)}
                    >
                        {'Posts'}
                    </div>
                    <div
                        className={classNames('menu-title', {active: selectedArea === BackstageArea.Users})}
                        onClick={() => setSelectedArea(BackstageArea.Users)}
                    >
                        {'Users'}
                    </div>
                </div>
            </div>
            <div className='content-container'>
                {activeArea}
            </div>
        </div>
    );
};

export default Backstage;
