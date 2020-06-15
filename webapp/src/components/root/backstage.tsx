// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import classNames from 'classnames';

import './backstage.scss';

interface Props {
    onBack: () => void;
    theme: Record<string, string>;
}

const Backstage = ({onBack}: Props): React.ReactElement<Props> => {
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
                    {/*<div className={classNames('menu-title', {active: selectedArea === BackstageArea.Dashboard})}>
                        {'Dashboard'}
                    </div>*/}
                    <div
                        className={classNames('menu-title')}
                        onClick={() => {}}
                    >
                        {'Posts'}
                    </div>
                    <div
                        className={classNames('menu-title')}
                        onClick={() => {}}
                    >
                        {'Users'}
                    </div>
                </div>
            </div>
            <div className='content-container'>
                {<div> data </div>}
            </div>
        </div>
    );
};

export default Backstage;
