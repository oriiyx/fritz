import {DataComponent} from '@/generated/definitions'
import {Input} from '@/components/Input'

interface CommonSettingsFieldsProps {
    component: DataComponent
    onFieldChange: (field: keyof DataComponent, value: any) => void
}

export function CommonSettingsFields({component, onFieldChange}: CommonSettingsFieldsProps) {
    return (
        <>
            {/* Basic Information */}
            <div>
                <h3 className="mb-4 text-sm font-semibold text-base-content/70">Basic Information</h3>
                <div className="space-y-4">
                    <div className="fieldset">
                        <label className="label">
                            <span className="label-text">Component Type</span>
                        </label>
                        <div className="badge badge-primary badge-sm">{component.type}</div>
                    </div>

                    <Input
                        label="Field Name (database column)"
                        value={component.name}
                        onChange={(e) => onFieldChange('name', e.target.value)}
                        fullWidth
                        disabled
                        className="font-mono"
                    />

                    <Input
                        label="Display Title"
                        value={component.title}
                        onChange={(e) => onFieldChange('title', e.target.value)}
                        fullWidth
                    />

                    <div className="fieldset">
                        <label className="label">
                            <span className="label-text">Database Type</span>
                        </label>
                        <div className="badge badge-neutral badge-sm">{component.dbtype}</div>
                    </div>
                </div>
            </div>

            <div className="divider"></div>

            {/* Flags */}
            <div>
                <h3 className="mb-4 text-sm font-semibold text-base-content/70">Field Properties</h3>
                <div className="space-y-3">
                    <div className="form-control fieldset">
                        <label className="label cursor-pointer justify-start gap-3">
                            <input
                                type="checkbox"
                                checked={component.mandatory}
                                onChange={(e) => onFieldChange('mandatory', e.target.checked)}
                                className="checkbox checkbox-primary"
                            />
                            <div>
                                <span className="label-text font-medium">Mandatory</span>
                                <p className="text-xs text-base-content/60">
                                    This field must have a value
                                </p>
                            </div>
                        </label>
                    </div>

                    <div className="form-control fieldset">
                        <label className="label cursor-pointer justify-start gap-3">
                            <input
                                type="checkbox"
                                checked={component.invisible}
                                onChange={(e) => onFieldChange('invisible', e.target.checked)}
                                className="checkbox checkbox-primary"
                            />
                            <div>
                                <span className="label-text font-medium">Invisible</span>
                                <p className="text-xs text-base-content/60">
                                    Hide this field from forms
                                </p>
                            </div>
                        </label>
                    </div>

                    <div className="form-control fieldset">
                        <label className="label cursor-pointer justify-start gap-3">
                            <input
                                type="checkbox"
                                checked={component.notEditable}
                                onChange={(e) => onFieldChange('notEditable', e.target.checked)}
                                className="checkbox checkbox-primary"
                            />
                            <div>
                                <span className="label-text font-medium">Not Editable</span>
                                <p className="text-xs text-base-content/60">
                                    Field is read-only
                                </p>
                            </div>
                        </label>
                    </div>
                </div>
            </div>
        </>
    )
}
