import React from 'react';

export default function Footer(): JSX.Element {
    return (
        <footer className="relative pt-8 pb-6 bg-teal-400">
            <div className="flex flex-wrap items-center justify-center md:justify-between">
                <div className="w-full px-4 mx-auto text-center md:w-4/12">
                    <div className="py-1 text-lg font-semibold text-blue-600">
                        Copyright © {new Date().getFullYear()} Avian Project{' '}
                        <a href="https://chingu.io/" className="text-lg text-blue-600 hover:text-blue-500">
                            Chingu
                        </a>
                    </div>
                </div>
            </div>
        </footer>
    );
}
