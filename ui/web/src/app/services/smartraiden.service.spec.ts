import { TestBed, inject } from '@angular/core/testing';

import { SmartRaidenService } from './smartraiden.service';

describe('SmartRaidenService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [SmartRaidenService]
    });
  });

  it('should ...', inject([SmartRaidenService], (service: SmartRaidenService) => {
    expect(service).toBeTruthy();
  }));
});
