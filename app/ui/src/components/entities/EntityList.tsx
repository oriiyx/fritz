import {Card, CardBody, CardTitle} from "@/components/Card.tsx";
import {CubeIcon} from "@heroicons/react/24/outline";
import {EntityDefinition} from "@/stores/entitiesDefinitionsStore.ts";
import {Loading} from "@/components/Loading.tsx";

interface EntityListProps {
    entityDefinitions: EntityDefinition[] | null
    selectedDefinition: EntityDefinition | null
    handleSelectDefinition: (definition: EntityDefinition) => void
}

export function EntityList({
                               entityDefinitions,
                               selectedDefinition,
                               handleSelectDefinition
                           }: EntityListProps) {

    if (!entityDefinitions) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loading variant="spinner" size="lg"/>
            </div>
        )
    }

    return (
        <Card>
            <CardBody>
                <CardTitle>Available Definitions</CardTitle>
                {entityDefinitions.length === 0 ? (
                    <div className="py-12 text-center text-base-content/50">
                        <CubeIcon className="mx-auto h-12 w-12 opacity-50"/>
                        <p className="mt-4">No entity definitions yet</p>
                        <p className="mt-1 text-sm">Create your first one to get started</p>
                    </div>
                ) : (
                    <div className="space-y-2">
                        {entityDefinitions.map((definition) => (
                            <button
                                key={definition.id}
                                onClick={() => handleSelectDefinition(definition)}
                                className={`w-full rounded-lg border-2 p-2 text-left transition-all hover:border-primary ${
                                    selectedDefinition?.id === definition.id
                                        ? 'border-primary bg-primary/10'
                                        : 'border-base-300'
                                }`}
                            >
                                <div className="flex items-start justify-between">
                                    <div className="flex-1">
                                        <h3 className="font-semibold">{definition.name}</h3>
                                        <p className="text-sm text-base-content/60">
                                            ID: {definition.id}
                                        </p>
                                    </div>
                                </div>
                            </button>
                        ))}
                    </div>
                )}
            </CardBody>
        </Card>
    );
}