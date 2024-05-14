## Discussion
**Scenario -** Option to delete a BG from catalog when it has been added to user lists
**Assumption -** On GetLIst() we always fetch each BG information to display other stuff
**Options:**
- If on BG lookup, it is not found, we can delete that entry from the List immediately
- - **Advantages -** Database space is optimized 
- - **Disadvantages -** May require constant purging of database entries which is not only heavier (more requests) but can also create unexpected issues. Includes more requests to different services.

- Catalog entries cannot be deleted and we flag the BG entry on the catalog as "discontinued". On BG lookup, we can retrieve that, displaying to the the user an error indicating that the BG no longer exists, allowing the entry to be removed from the list.
- - **Advantages -** Append only is less prone to issues and overall faster since no deletes are made on the fly
- - **Disadvantages -** Requires more verifications on the FE
  
**Choice -** Option 2