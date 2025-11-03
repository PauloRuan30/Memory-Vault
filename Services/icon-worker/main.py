import os
import json
import time
import redis
import psycopg2
from minio import Minio
from processor import generate_texture

def main():
    # Load environment variables
    db_host = os.getenv('DB_HOST', 'localhost')
    db_port = os.getenv('DB_PORT', '5432')
    db_user = os.getenv('DB_USER', 'postgres')
    db_password = os.getenv('DB_PASSWORD', 'postgres')
    db_name = os.getenv('DB_NAME', 'memory_vault')
    redis_addr = os.getenv('REDIS_ADDR', 'localhost:6379')
    s3_endpoint = os.getenv('S3_ENDPOINT', 'http://localhost:9000')
    s3_access_key = os.getenv('S3_ACCESS_KEY', 'minioadmin')
    s3_secret_key = os.getenv('S3_SECRET_KEY', 'minioadmin')
    s3_bucket = os.getenv('S3_BUCKET', 'memory-vault')

    # Initialize connections
    redis_client = redis.Redis(host=redis_addr.split(':')[0], port=int(redis_addr.split(':')[1]), decode_responses=True)
    db_conn = psycopg2.connect(
        host=db_host,
        port=db_port,
        user=db_user,
        password=db_password,
        database=db_name
    )
    s3_client = Minio(
        s3_endpoint.replace('http://', '').replace('https://', ''),
        access_key=s3_access_key,
        secret_key=s3_secret_key,
        secure=False
    )

    print("Icon worker started, waiting for jobs...")

    while True:
        try:
            # Blocking pop from job queue
            result = redis_client.brpop(['icon_jobs'], timeout=5)
            if result:
                _, job_data = result
                job = json.loads(job_data)
                file_id = job['file_id']
                
                print(f"Processing job for file_id: {file_id}")
                
                # Update status to PROCESSING
                cur = db_conn.cursor()
                cur.execute(
                    "UPDATE files SET processing_status = 'PROCESSING' WHERE id = %s",
                    (file_id,)
                )
                db_conn.commit()
                cur.close()
                
                # Download file from S3
                cur = db_conn.cursor()
                cur.execute(
                    "SELECT s3_path FROM files WHERE id = %s",
                    (file_id,)
                )
                row = cur.fetchone()
                if not row:
                    print(f"File {file_id} not found in database")
                    continue
                
                s3_path = row[0]
                temp_filename = f"/tmp/{file_id}_temp"
                
                s3_client.fget_object(s3_bucket, s3_path, temp_filename)
                cur.close()
                
                # Generate texture
                texture_filename = f"{file_id}_texture.png"
                texture_path = f"icons/{texture_filename}"
                generate_texture(temp_filename, f"/tmp/{texture_filename}")
                
                # Upload texture to S3
                s3_client.fput_object(
                    s3_bucket,
                    texture_path,
                    f"/tmp/{texture_filename}",
                    content_type='image/png'
                )
                
                # Update database
                cur = db_conn.cursor()
                cur.execute(
                    "UPDATE files SET processing_status = 'COMPLETE', texture_path = %s WHERE id = %s",
                    (texture_path, file_id)
                )
                db_conn.commit()
                
                # Publish update to Redis channel for real-time notifications
                redis_client.publish('file_updates', json.dumps({
                    'file_id': file_id,
                    'texture_path': texture_path
                }))
                
                cur.close()
                
                # Cleanup temp files
                os.remove(temp_filename)
                os.remove(f"/tmp/{texture_filename}")
                
                print(f"Successfully processed file {file_id}")
                
        except Exception as e:
            print(f"Error processing job: {e}")
            time.sleep(5)  # Wait before retrying

if __name__ == "__main__":
    main()