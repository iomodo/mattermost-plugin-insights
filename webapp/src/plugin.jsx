import React from 'react';

import en from '../i18n/en.json';
import {id as pluginId} from './manifest';
import { openRootModal } from './actions';
import Root from './components/root';
import reducer from './reducer';


function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    }
    return {};
}

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerRootComponent(Root);

        registry.registerMainMenuAction(
            'Workplace Insights',
            () => store.dispatch(openRootModal())
        );
        registry.registerReducer(reducer);

        registry.registerTranslations(getTranslations);
    }
}