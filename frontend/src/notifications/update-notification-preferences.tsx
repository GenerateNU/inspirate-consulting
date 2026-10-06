import { useEffect, useState } from 'react';
import { useGetNotificationPreferences, useUpdateNotificationPreferences } from '../api/endpoints';

export default function UpdateNotificationPreferences() {
    const { data: response, error, isLoading } = useGetNotificationPreferences();
    const { trigger, isMutating } = useUpdateNotificationPreferences();

    const [emailEnabled, setEmailEnabled] = useState(true);
    const [weeklySummaryEnabled, setWeeklySummaryEnabled] = useState(true);
    const [dueDateNotificationsEnabled, setDueDateNotificationsEnabled] = useState(true);
    const [daysBeforeDue, setDaysBeforeDue] = useState(1);
    const [notifyPastDue, setNotifyPastDue] = useState(true);
    const [status, setStatus] = useState('');

    // Pre-fill the form with GET request (current preferences).
    useEffect(() => {
        if (response?.status === 200) {
            const prefs = response.data;
            setEmailEnabled(prefs.email_enabled);
            setWeeklySummaryEnabled(prefs.weekly_summary_enabled);
            setDueDateNotificationsEnabled(prefs.due_date_notifications_enabled);
            setDaysBeforeDue(prefs.days_before_due);
            setNotifyPastDue(prefs.notify_past_due);
        }
    }, [response]);

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const result = await trigger({
                email_enabled: emailEnabled,
                weekly_summary_enabled: weeklySummaryEnabled,
                due_date_notifications_enabled: dueDateNotificationsEnabled,
                days_before_due: daysBeforeDue,
                notify_past_due: notifyPastDue,
            });

            if (result.status === 200) {
                setStatus('Saved.');
            } else {
                setStatus('Could not save preferences.');
            }
        } catch (err) {
            console.error(err);
            setStatus('Something went wrong.');
        }
    };

    if (isLoading) {
        return (
            <div className="form-container">
                <h2>Notification Preferences</h2>
                <p>Loading...</p>
            </div>
        );
    }

    if (error || !response || response.status !== 200) {
        return (
            <div className="form-container">
                <h2>Notification Preferences</h2>
                <p>Something went wrong.</p>
            </div>
        );
    }

    return (
        <div className="form-container">
            <h2>Notification Preferences</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>
                        <input
                            type="checkbox"
                            checked={emailEnabled}
                            onChange={(e) => setEmailEnabled(e.target.checked)}
                        />
                        Email notifications enabled
                    </label>
                </div>

                <div className="form-group">
                    <label>
                        <input
                            type="checkbox"
                            checked={weeklySummaryEnabled}
                            onChange={(e) => setWeeklySummaryEnabled(e.target.checked)}
                        />
                        Weekly summary of outstanding assignments
                    </label>
                </div>

                <div className="form-group">
                    <label>
                        <input
                            type="checkbox"
                            checked={dueDateNotificationsEnabled}
                            onChange={(e) => setDueDateNotificationsEnabled(e.target.checked)}
                        />
                        Reminder before an assignment is due
                    </label>
                </div>

                <div className="form-group">
                    <label>Days before due date to send the reminder:</label>
                    <input
                        type="number"
                        min={0}
                        value={daysBeforeDue}
                        onChange={(e) => setDaysBeforeDue(Number(e.target.value))}
                        disabled={!dueDateNotificationsEnabled}
                    />
                </div>

                <div className="form-group">
                    <label>
                        <input
                            type="checkbox"
                            checked={notifyPastDue}
                            onChange={(e) => setNotifyPastDue(e.target.checked)}
                        />
                        Notify when an assignment becomes past due
                    </label>
                </div>

                <button type="submit" disabled={isMutating}>
                    {isMutating ? 'Saving...' : 'Save Preferences'}
                </button>
            </form>
            {status && <p>{status}</p>}
        </div>
    );
}