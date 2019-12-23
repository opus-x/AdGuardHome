import React from 'react';
import PropTypes from 'prop-types';
import { Field, reduxForm } from 'redux-form';
import { Trans, withNamespaces } from 'react-i18next';
import flow from 'lodash/flow';

import { renderSelectField, renderRadioField, toNumber } from '../../../helpers/form';
import { QUERY_LOG_INTERVALS_DAYS } from '../../../helpers/constants';

const getIntervalFields = (processing, t, toNumber) =>
    QUERY_LOG_INTERVALS_DAYS.map((interval) => {
        const title =
            interval === 1 ? t('interval_24_hour') : t('interval_days', { count: interval });

        return (
            <Field
                key={interval}
                name="interval"
                type="radio"
                component={renderRadioField}
                value={interval}
                placeholder={title}
                normalize={toNumber}
                disabled={processing}
            />
        );
    });

const Form = (props) => {
    const {
        handleSubmit,
        submitting,
        invalid,
        processing,
        processingClear,
        handleClear,
        t,
        interval,
        enabled,
    } = props;

    return (
        <form onSubmit={handleSubmit}>
            <div className="form__group form__group--settings">
                <Field
                    name="enabled"
                    type="checkbox"
                    component={renderSelectField}
                    placeholder={t('query_log_enable')}
                    disabled={processing}
                />
            </div>
            {(enabled && interval !== 1) &&
            <div className="form__group form__group--settings">
                <Field
                    name="displayDate"
                    type="checkbox"
                    component={renderSelectField}
                    placeholder={t('display_date')}
                    disabled={processing}
                    subtitle={t('query_log_display_date')}
                />
                <div data-tooltip={t('requests_not_made_today')} className="tooltip-custom tooltip-custom--narrow tooltip-custom--narrow tooltip-custom--lower"/>
            </div>}
            <label className="form__label">
                <Trans>query_log_retention</Trans>
            </label>
            <div className="form__group form__group--settings">
                <div className="custom-controls-stacked">
                    {getIntervalFields(processing, t, toNumber)}
                </div>
            </div>
            <div className="mt-5">
                <button
                    type="submit"
                    className="btn btn-success btn-standard btn-large"
                    disabled={submitting || invalid || processing}
                >
                    <Trans>save_btn</Trans>
                </button>
                <button
                    type="button"
                    className="btn btn-outline-secondary btn-standard ml-5"
                    onClick={() => handleClear()}
                    disabled={processingClear}
                >
                    <Trans>query_log_clear</Trans>
                </button>
            </div>
        </form>
    );
};

Form.propTypes = {
    handleSubmit: PropTypes.func.isRequired,
    handleClear: PropTypes.func.isRequired,
    submitting: PropTypes.bool.isRequired,
    invalid: PropTypes.bool.isRequired,
    processing: PropTypes.bool.isRequired,
    processingClear: PropTypes.bool.isRequired,
    t: PropTypes.func.isRequired,
    interval: PropTypes.number.isRequired,
    enabled: PropTypes.bool.isRequired,
};

export default flow([
    withNamespaces(),
    reduxForm({
        form: 'logConfigForm',
    }),
])(Form);
