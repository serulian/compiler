$module('async', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('5f7de1e8', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.SomeGenerator = function () {
    var $result;
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(false, $g.________testlib.basictypes.Boolean));
            $current = 1;
            return;

          case 1:
            $promise.translate($g.async.DoSomethingAsync()).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($yield, $yieldin, $reject, $done);
              return;
            }).catch(function (err) {
              throw err;
            });
            return;

          case 2:
            $yield($result);
            $current = 3;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue, true, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var v;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            v = null;
            $current = 1;
            continue localasyncloop;

          case 1:
            $promise.maybe($g.async.SomeGenerator()).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $temp1 = $result;
            $current = 3;
            continue localasyncloop;

          case 3:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 5;
              continue localasyncloop;
            } else {
              $current = 6;
              continue localasyncloop;
            }
            break;

          case 5:
            v = value;
            $current = 3;
            continue localasyncloop;

          case 6:
            $resolve(v);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
