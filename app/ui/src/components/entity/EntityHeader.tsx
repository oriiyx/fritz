import {Entity} from "@/services/entitiesService.ts";
import {Button} from "@/components/Button.tsx";
import {ArrowLeftIcon} from "@heroicons/react/24/outline";
import {useNavigate} from "@tanstack/react-router";
import {EntityDefinition} from "@/generated/definitions";

interface EntityHeaderProps {
    entity: Entity
    definition: EntityDefinition
}

export function EntityHeader({entity, definition}: EntityHeaderProps) {
    const navigate = useNavigate()

    return (<div className="mb-6">
        <div className="flex items-center gap-4 mb-4">
            <Button
                variant="ghost"
                size="sm"
                iconLeft={<ArrowLeftIcon className="h-4 w-4"/>}
                onClick={() => navigate({to: '/'})}
            >
                Back
            </Button>
        </div>
        <div className="flex items-center gap-3">
            <div>
                <h1 className="text-2xl font-bold">{entity.o_key}</h1>
                <p className="text-sm text-base-content/70">
                    {entity.o_path} • {definition.name}
                </p>
            </div>
            <div className="flex gap-2 ml-auto">
                {entity.published ? (
                    <span className="badge badge-success">Published</span>
                ) : (
                    <span className="badge badge-warning">Draft</span>
                )}
                {!entity.has_data && (
                    <span className="badge badge-info">No data yet</span>
                )}
            </div>
        </div>
    </div>)
}
