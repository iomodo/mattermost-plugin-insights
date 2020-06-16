import React from 'react';

import en from '../i18n/en.json';
import { setBackstageModal } from './actions';
import reducer from './reducer';
import BackstageModal from './components/backstage/backstage_modal';
import {BackstageArea} from './types/backstage';


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
        registry.registerRootComponent(BackstageModal);

        registry.registerMainMenuAction(
            'Workplace Insights',
            () => store.dispatch(setBackstageModal(true, BackstageArea.Users))
        );
        registry.registerReducer(reducer);

        registry.registerTranslations(getTranslations);
    }
}