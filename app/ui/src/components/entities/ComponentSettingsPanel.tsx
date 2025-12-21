import {DataComponent} from '@/generated/definitions'
import {Input} from '@/components/Input'
import {Card, CardBody, CardTitle} from '@/components/Card'

interface Props {
    component: DataComponent | null
    onComponentUpdate: (updatedComponent: DataComponent) => void
}

export function ComponentSettingsPanel({component, onComponentUpdate}: Props) {
    if (!component) {
        return (
            <Card>
                <CardBody>
                    <div className="flex h-96 items-center justify-center text-base-content/50">
                        <div className="text-center">
                            <p className="text-lg">No component selected</p>
                            <p className="mt-2 text-sm">Select a component from the tree to edit its settings</p>
                        </div>
                    </div>
                </CardBody>
            </Card>
        )
    }

    const handleFieldChange = (field: keyof DataComponent, value: any) => {
        onComponentUpdate({
            ...component,
            [field]: value,
        })
    }

    const handleSettingsChange = (settingKey: string, value: any) => {
        onComponentUpdate({
            ...component,
            settings: {
                ...component.settings,
                [settingKey]: value,
            },
        })
    }

    return (
        <Card>
            <CardBody>
                <CardTitle>Component Settings</CardTitle>

                <div className="mt-6 space-y-6">
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
                                onChange={(e) => handleFieldChange('name', e.target.value)}
                                fullWidth
                                disabled
                                className="font-mono"
                            />

                            <Input
                                label="Display Title"
                                value={component.title}
                                onChange={(e) => handleFieldChange('title', e.target.value)}
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
                                        onChange={(e) => handleFieldChange('mandatory', e.target.checked)}
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
                                        onChange={(e) => handleFieldChange('invisible', e.target.checked)}
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
                                        onChange={(e) => handleFieldChange('notEditable', e.target.checked)}
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

                    <div className="divider"></div>

                    {/* Type-Specific Settings */}
                    <div>
                        <h3 className="mb-4 text-sm font-semibold text-base-content/70">
                            {component.type.charAt(0).toUpperCase() + component.type.slice(1)} Settings
                        </h3>

                        {component.type === 'input' && (
                            <div className="space-y-4">
                                <Input
                                    label="Default Value"
                                    value={(component.settings as any).defaultValue || ''}
                                    onChange={(e) => handleSettingsChange('defaultValue', e.target.value)}
                                    fullWidth
                                />
                                <Input
                                    label="Column Length"
                                    type="number"
                                    value={(component.settings as any).columnLength || ''}
                                    onChange={(e) =>
                                        handleSettingsChange(
                                            'columnLength',
                                            e.target.value ? parseInt(e.target.value) : null
                                        )
                                    }
                                    fullWidth
                                />
                                <Input
                                    label="Regex Validation"
                                    value={(component.settings as any).regexValidation || ''}
                                    onChange={(e) => handleSettingsChange('regexValidation', e.target.value)}
                                    fullWidth
                                    className="font-mono text-sm"
                                />
                            </div>
                        )}

                        {component.type === 'integer' && (
                            <div className="space-y-4">
                                <Input
                                    label="Default Value"
                                    type="number"
                                    value={(component.settings as any).defaultValue || ''}
                                    onChange={(e) =>
                                        handleSettingsChange(
                                            'defaultValue',
                                            e.target.value ? parseInt(e.target.value) : null
                                        )
                                    }
                                    fullWidth
                                />
                                <Input
                                    label="Minimum Value"
                                    type="number"
                                    value={(component.settings as any).minValue || ''}
                                    onChange={(e) =>
                                        handleSettingsChange(
                                            'minValue',
                                            e.target.value ? parseInt(e.target.value) : null
                                        )
                                    }
                                    fullWidth
                                />
                                <Input
                                    label="Maximum Value"
                                    type="number"
                                    value={(component.settings as any).maxValue || ''}
                                    onChange={(e) =>
                                        handleSettingsChange(
                                            'maxValue',
                                            e.target.value ? parseInt(e.target.value) : null
                                        )
                                    }
                                    fullWidth
                                />
                                <div className="form-control fieldset">
                                    <label className="label cursor-pointer justify-start gap-3">
                                        <input
                                            type="checkbox"
                                            checked={(component.settings as any).unsigned || false}
                                            onChange={(e) =>
                                                handleSettingsChange('unsigned', e.target.checked)
                                            }
                                            className="checkbox checkbox-primary"
                                        />
                                        <span className="label-text">Unsigned (positive numbers only)</span>
                                    </label>
                                </div>
                            </div>
                        )}

                        {component.type === 'date' && (
                            <div className="space-y-4">
                                <Input
                                    label="Default Value"
                                    type="date"
                                    value={(component.settings as any).defaultValue || ''}
                                    onChange={(e) => handleSettingsChange('defaultValue', e.target.value)}
                                    fullWidth
                                />
                            </div>
                        )}
                    </div>
                </div>
            </CardBody>
        </Card>
    )
}