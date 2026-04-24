-- Seed default categories
INSERT INTO categories (id, name, slug, description, created_at, updated_at) VALUES
(gen_random_uuid(), 'Technology', 'technology', 'All things related to new tech, software, and hardware.', NOW(), NOW()),
(gen_random_uuid(), 'Business', 'business', 'Corporate talk, start-ups, and market analysis.', NOW(), NOW()),
(gen_random_uuid(), 'Lifestyle', 'lifestyle', 'Health, wellness, and daily living topics.', NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

-- Seed default blogs using the new categories and the existing 'admin@gns.com' user
-- Note: 'admin@gns.com' is seeded from rbac_seeder.sql
INSERT INTO blogs (id, title, slug, content, excerpt, cover_image_url, status, category_id, author_id, published_at, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'The Future of AI in Development',
    'future-of-ai-in-development',
    '<p>Artificial intelligence is rapidly changing how we write code. This post explores the immediate future of LLMs in day-to-day development.</p>',
    'A quick dive into LLMs in coding.',
    'https://images.unsplash.com/photo-1549692520-acc6669e2f0c?q=80&w=1500',
    'published',
    c.id,
    u.id,
    NOW(),
    NOW(),
    NOW()
FROM categories c
CROSS JOIN users u
WHERE c.slug = 'technology' AND u.email = 'admin@gns.com'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO blogs (id, title, slug, content, excerpt, cover_image_url, status, category_id, author_id, published_at, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'Starting Your First Tech Company',
    'starting-your-first-tech-company',
    '<p>Building a startup is hard. Here are 10 things you must know before taking the leap.</p>',
    '10 things about starting a company.',
    'https://images.unsplash.com/photo-1556761175-5973dc0f32d7?q=80&w=1500',
    'draft',
    c.id,
    u.id,
    NULL,
    NOW(),
    NOW()
FROM categories c
CROSS JOIN users u
WHERE c.slug = 'business' AND u.email = 'admin@gns.com'
ON CONFLICT (slug) DO NOTHING;
