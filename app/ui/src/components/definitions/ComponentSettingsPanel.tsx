import {DataComponent} from '@/generated/definitions'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {CommonSettingsFields, isFloat4Component, isFloat8Component} from '@/components/definitions/settings'
import {InputSettings} from './settings/InputSettings'
import {IntegerSettings} from './settings/IntegerSettings'
import {DateSettings} from './settings/DateSettings'
import {isDateComponent, isInputComponent, isIntegerComponent,} from '@/components/definitions/settings'
import {FloatSettings} from "@/components/definitions/settings/FloatSettings.tsx";

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

    const handleFieldChange = (field: keyof DataComponent, value: unknown) => {
        onComponentUpdate({
            ...component,
            [field]: value,
        })
    }

    const handleSettingsChange = (settingKey: string, value: unknown) => {
        onComponentUpdate({
            ...component,
            settings: {
                ...component.settings,
                [settingKey]: value,
            },
        })
    }

    // Render type-specific settings based on component type
    const renderTypeSpecificSettings = () => {
        if (isInputComponent(component)) {
            return <InputSettings component={component} onSettingsChange={handleSettingsChange}/>
        }

        if (isIntegerComponent(component)) {
            return <IntegerSettings component={component} onSettingsChange={handleSettingsChange}/>
        }

        if (isFloat4Component(component) || isFloat8Component(component)) {
            return <FloatSettings component={component} onSettingsChange={handleSettingsChange}/>
        }

        if (isDateComponent(component)) {
            return <DateSettings component={component} onSettingsChange={handleSettingsChange}/>
        }

        // Fallback for component types without specific settings
        return (
            <div className="text-center text-base-content/50 py-8">
                <p className="text-sm">No additional settings available for this component type</p>
            </div>
        )
    }

    return (
        <Card>
            <CardBody>
                <CardTitle>Component Settings</CardTitle>

                <div className="mt-6 space-y-6">
                    {/* Common fields for all component types */}
                    <CommonSettingsFields component={component} onFieldChange={handleFieldChange}/>

                    <div className="divider"></div>

                    {/* Type-Specific Settings */}
                    <div>
                        <h3 className="mb-4 text-sm font-semibold text-base-content/70">
                            {component.type.charAt(0).toUpperCase() + component.type.slice(1)} Settings
                        </h3>

                        {renderTypeSpecificSettings()}
                    </div>
                </div>
            </CardBody>
        </Card>
    )
}