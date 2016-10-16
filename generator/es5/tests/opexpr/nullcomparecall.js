$module('nullcomparecall', function () {
  var $static = this;
  this.$class('SomeError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box('huh?', $g.____testlib.basictypes.String));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Message', 3, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.nullcomparecall.SomeError).$typeref()]);
    };
  });

  $static.RejectNow = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.nullcomparecall.SomeError.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $reject($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var thing;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            thing = $t.box(true, $g.____testlib.basictypes.Boolean);
            $promise.resolve(thing).then(function ($result0) {
              return ($promise.shortcircuit($result0, null) || $g.nullcomparecall.RejectNow()).then(function ($result1) {
                $result = $t.nullcompare($result0, $result1);
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
