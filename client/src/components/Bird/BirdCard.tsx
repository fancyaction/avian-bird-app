import React from 'react';
import { Card, CardContent } from 'utils';
import ExampleImage from 'assets/examplebird.jpg';
import BirdImage from './BirdImage';
import BirdDetails, { BirdDetailsProps } from './BirdDetails';

const PlaceholderBird: React.FC = () => (
    <div className="flex flex-col items-center justify-center p-4">
        <BirdImage imageUrl={ExampleImage} imageAlt="Placeholder bird" />
        <h1 className="py-2 text-2xl font-bold text-green-800">
            Start by uploading an image via the "Send Prediction" button above
        </h1>
    </div>
);

const BirdCard: React.FC<BirdDetailsProps> = ({ imageUrl, imageAlt = 'Photo of bird', loading, record, children }) => {
    return (
        <Card color="bg-teal-400">
            <CardContent>
                {!record && !loading ? (
                    <PlaceholderBird />
                ) : (
                    <BirdDetails imageUrl={imageUrl} imageAlt={imageAlt} record={record} loading={loading} />
                )}
            </CardContent>
            {children}
        </Card>
    );
};

export default BirdCard;
