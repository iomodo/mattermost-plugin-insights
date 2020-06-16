import {SetBackstageModalSettings, SET_BACKSTAGE_MODAL_SETTINGS} from './types/actions'
import {BackstageArea} from './types/backstage';

export function setBackstageModal(open: boolean, selectedArea?: BackstageArea): SetBackstageModalSettings {
    return {
        type: SET_BACKSTAGE_MODAL_SETTINGS,
        open,
        selectedArea,
    };
}