import React from 'react';
import { useAuthStore, UserRole } from '../store/authStore';

interface ShowForRoleProps {
    allowedRoles: UserRole[];
    children: React.ReactNode;
    fallback?: React.ReactNode;
}

/**
 * Conditional rendering component based on user role.
 * Only renders children if user has one of the allowed roles.
 */
export const ShowForRole: React.FC<ShowForRoleProps> = ({
    allowedRoles,
    children,
    fallback = null,
}) => {
    const { hasRole } = useAuthStore();

    if (hasRole(allowedRoles)) {
        return <>{children}</>;
    }

    return <>{fallback}</>;
};

export default ShowForRole;
