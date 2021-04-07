import React, { useCallback, useState } from 'react';
import parse from 'html-react-parser';
import SpeciesInfoTable from './SpeciesInfoTable';
import { useBreakpoints, Button } from 'utils';
import axios from 'axios';
import BirdImage, { BirdImageProps } from './BirdImage';

const EBirdRedirectButton: React.FC<{ name: string }> = ({ name }) => {
    const redirectToEBirdUrl = useCallback(() => {
        const config = {
            params: {
                locale: 'en_US',
                cat: 'species',
                limit: 150,
                key: 'jfekjedvescr',
                q: name,
            },
        };
        axios
            .get(`https://api.ebird.org/v2/ref/taxon/find`, config)
            .then((response) => {
                const speciesCode = response.data[0].code;
                window.open(`https://ebird.org/species/${speciesCode}`, '_blank');
            })
            .catch((error) => {
                console.log(error);
            });
    }, [name]);

    return <Button onClick={redirectToEBirdUrl} label="Read More" />;
};

type BirdRecord = {
    name: string;
    species_info: {
        synonyms: string;
        otherCommonNames: string[];
        kingdom: string;
        phylum: string;
        taxclass: string;
        taxorder: string;
        family: string;
        genus: string;
        informalTaxonomy: string;
        taxonomicComments: string;
    };
};
interface ViewProps extends BirdImageProps {
    record: BirdRecord;
}

const MobileView: React.FC<ViewProps> = ({ record, imageUrl, imageAlt }) => {
    const [showMore, setShowMore] = useState<boolean>(false);
    const { name, species_info: speciesInfo } = record;
    const hasSpeciesInfo = 0 < Object.keys(speciesInfo).length;
    const { taxonomicComments: description, ...otherInfo } = speciesInfo;

    const Content = useCallback(
        () =>
            showMore && hasSpeciesInfo ? (
                <SpeciesInfoTable {...otherInfo} />
            ) : (
                <>
                    <h1 className="py-2 text-2xl font-bold text-green-800">{name}</h1>
                    <p className="mt-4 text-sm text-black break-all md:break-all ">
                        {hasSpeciesInfo ? parse(description) : ''}
                    </p>
                </>
            ),
        [description, hasSpeciesInfo, name, otherInfo, showMore]
    );

    return (
        <div className="p-2">
            {imageUrl && <BirdImage imageUrl={imageUrl} imageAlt={imageAlt} />}
            <Content />
            {hasSpeciesInfo && (
                <Button onClick={() => setShowMore(!showMore)} label={showMore ? 'Hide Details' : 'Show Details'} />
            )}
            <EBirdRedirectButton name={name} />
        </div>
    );
};

const DesktopView: React.FC<ViewProps> = ({ record, imageUrl, imageAlt }) => {
    const { name, species_info: speciesInfo } = record;
    const hasSpeciesInfo = 0 < Object.keys(speciesInfo).length;
    const { taxonomicComments: description, ...otherInfo } = speciesInfo;

    return (
        <>
            <div className="flex flex-col items-center justify-center p-4">
                {imageUrl && <BirdImage imageUrl={imageUrl} imageAlt={imageAlt} />}
                <h1 className="py-2 text-2xl font-bold text-green-800">{name}</h1>
                <p className="mt-4 text-sm text-black break-all md:break-all ">
                    {hasSpeciesInfo ? parse(description) : ''}
                </p>
                <EBirdRedirectButton name={name} />
            </div>
            <div className="p-4">{hasSpeciesInfo && <SpeciesInfoTable {...otherInfo} />}</div>
        </>
    );
};

export interface BirdDetailsProps extends ViewProps {
    loading: boolean;
}

const BirdDetails: React.FC<BirdDetailsProps> = ({ loading, ...props }) => {
    const { isXs, isSm } = useBreakpoints();
    const isMobile = isXs || isSm;

    if (loading) {
        return <span>Loading...</span>;
    }

    if (isMobile) {
        return <MobileView {...props} />;
    }

    return <DesktopView {...props} />;
};

export default BirdDetails;
