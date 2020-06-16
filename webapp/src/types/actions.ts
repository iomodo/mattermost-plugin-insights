import {id as pluginId} from '../manifest';
import {BackstageArea} from './backstage';

export const SET_BACKSTAGE_MODAL_SETTINGS = pluginId + '_set_backstage_modal_settings';

export interface SetBackstageModalSettings {
    type: typeof SET_BACKSTAGE_MODAL_SETTINGS;
    open: boolean;
    selectedArea?: BackstageArea;
}
