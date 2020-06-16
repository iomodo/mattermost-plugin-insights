import {combineReducers} from 'redux';
import {SET_BACKSTAGE_MODAL_SETTINGS} from './types/actions';


function backstageModal(state = {open: false, selectedArea: 0}, action) {
    switch (action.type) {
    case SET_BACKSTAGE_MODAL_SETTINGS:
        return {
            open: Boolean(action.open),
            selectedArea: action.selectedArea,
        };
    default:
        return state;
    }
}

export default combineReducers({
    backstageModal,
});