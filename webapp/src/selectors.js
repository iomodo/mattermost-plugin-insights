import {id as pluginId} from './manifest';


const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const backstageModal = (state) => getPluginState(state).backstageModal;
