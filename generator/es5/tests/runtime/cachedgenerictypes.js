$module('cachedgenerictypes', function () {
  var $static = this;
  this.$class('d23601e8', 'Boolean', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var firstMapping;
    var secondMapping;
    firstMapping = $g.________testlib.basictypes.Mapping($g.cachedgenerictypes.Boolean).Empty();
    secondMapping = $g.________testlib.basictypes.Mapping($g.________testlib.basictypes.Boolean).overObject((function () {
      var obj = {
      };
      obj["somekey"] = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      return obj;
    })());
    return $t.syncnullcompare(secondMapping.$index($t.fastbox('somekey', $g.________testlib.basictypes.String)), function () {
      return $t.fastbox(false, $g.________testlib.basictypes.Boolean);
    });
  };
});
