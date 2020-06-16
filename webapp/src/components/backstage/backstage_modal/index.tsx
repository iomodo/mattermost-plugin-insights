import {bindActionCreators, Dispatch} from 'redux';

import {connect} from 'react-redux';

import {setBackstageModal} from '../../../actions';
import {backstageModal} from '../../../selectors';
import {BackstageArea} from '../../../types/backstage';

import BackstageModal from './backstage_modal';

const mapStateToProps = (state: object): object => {
    return {
        show: backstageModal(state).open,
        selectedArea: backstageModal(state).selectedArea,
    };
};

const mapDispatchToProps = (dispatch: Dispatch): object =>
    bindActionCreators({
        close: () => setBackstageModal(false),
        setSelectedArea: (selectedArea: BackstageArea) => setBackstageModal(true, selectedArea),
    }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(BackstageModal);
