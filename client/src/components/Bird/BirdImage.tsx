import React from 'react';

export interface BirdImageProps {
    imageUrl: string;
    imageAlt: string;
}

const BirdImage: React.FC<BirdImageProps> = ({ imageUrl, imageAlt }) => (
    <div className="flex justify-center">
        <img
            src={imageUrl}
            alt={imageAlt}
            className="h-auto max-w-full align-middle border-none rounded shadow-lg max-h-72 md:max-h-40"
        />
    </div>
);

export default BirdImage;
