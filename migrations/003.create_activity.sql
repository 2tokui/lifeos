-- up
CREATE TABLE "Activity" (
	"id"	INTEGER NOT NULL UNIQUE,
	"task_id"	INTEGER NOT NULL,
	"mood_id"	INTEGER NOT NULL,
	"start_time"	INTEGER NOT NULL,
	"end_time"	INTEGER NOT NULL,
	"memo"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("mood_id") REFERENCES "",
	FOREIGN KEY("task_id") REFERENCES "Task"("id")
);
-- 
DROP TABLE "Activity";
